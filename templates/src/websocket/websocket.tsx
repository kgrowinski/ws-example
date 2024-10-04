import { createContext, FC, ReactNode, useCallback, useContext, useEffect, useRef, useState } from 'react';

import { v4 } from 'uuid';
import { ColorsContext } from '../colors/colors';

interface WebsocketContextInterface {
  isConnected: boolean;
  connection: WebSocket | null;
  handleSendWSMessage?: (type: string, payload?: any) => void;
}

export const WebsocketContext = createContext<WebsocketContextInterface>({
  isConnected: false,
  connection: null,
});

interface WebsocketContextProviderProps {
  children: ReactNode;
}

export const WS_ACTIONS = {
  WS_INIT: 'INIT_CONNECTION',
  WS_SET_COLOR: 'SET_COLOR',
  WS_NEW_COLOR: 'NEW_COLOR',
  WS_ERROR: 'ERROR',
};

type WSErrorPayload = {
  appDomain: string;
  errorCode: number;
}

type WSFileUploaded = {
  type: string;
  identifier: string;
};

type WSSetColorPayload = string;

type WSMessage = string;
type WSParsedMessage = {
  action: (typeof WS_ACTIONS)[keyof typeof WS_ACTIONS];
  payload:
    | string
    | WSErrorPayload
    | WSFileUploaded;
};

const RECONNECT_INTERVAL = 5000; // 5 seconds

export const WebsocketContextProvider: FC<WebsocketContextProviderProps> = ({ children }) => {
  const {setCurrentColor} = useContext(ColorsContext);
  const [isConnected, setIsConnected] = useState(false);

  const socketRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<any | null>(null);

  const clearReconnectTimeout = () => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }
  };

  const handleWSError = (err: WSErrorPayload) => {
    console.log(`${err.appDomain}.${err.errorCode}`);
  };

  const handleSendWSMessage = useCallback(async (type: string, payload?: object) => {
    const idToken = v4()
      socketRef.current.send(JSON.stringify({ Authorization: idToken, action: type, payload }));
  }, []);

  const handleWebsocketMessage = useCallback(
    (e: MessageEvent<WSMessage>) => {
      const parsedData: WSParsedMessage = JSON.parse(e.data);

      console.log('parsedData -', parsedData);
      switch (parsedData.action) {
        case WS_ACTIONS.WS_ERROR:
          handleWSError(parsedData.payload as WSErrorPayload);
          break;

        case WS_ACTIONS.WS_NEW_COLOR:
          setCurrentColor(parsedData.payload as string);
          break;

        default:
          console.warn(`Unhandled WebSocket action: ${parsedData.action}`);
      }
    },
    [ handleSendWSMessage],
  );

  const handleWebsocketConnection = useCallback(() => {
    clearReconnectTimeout();
    socketRef.current = new WebSocket(`ws://localhost:8080/ws/v1/websocket`);

    socketRef.current.onmessage = handleWebsocketMessage;

    socketRef.current.onopen = () => {
      handleSendWSMessage(WS_ACTIONS.WS_INIT);
      setIsConnected(true);
    };

    socketRef.current.onerror = (event: Event) => {
      console.error('WebSocket error:', event);
      setIsConnected(false);

      clearReconnectTimeout();
      reconnectTimeoutRef.current = setTimeout(() => {
        handleWebsocketConnection();
      }, RECONNECT_INTERVAL);
    };

    socketRef.current.onclose = () => {
      console.log('WebSocket connection closed');
      setIsConnected(false);

      clearReconnectTimeout();
      reconnectTimeoutRef.current = setTimeout(() => {
        handleWebsocketConnection();
      }, RECONNECT_INTERVAL);
    };
  }, [handleWebsocketMessage, handleSendWSMessage]);

  useEffect(() => {
    if ( !socketRef.current) {
      handleWebsocketConnection();
    }

    return () => {
      socketRef.current?.close();
      clearReconnectTimeout();
    };
  }, [handleWebsocketConnection]);

  return (
    <WebsocketContext.Provider
      value={{
        isConnected,
        connection: socketRef.current,
        handleSendWSMessage,
      }}
    >
      {children}
    </WebsocketContext.Provider>
  );
};
