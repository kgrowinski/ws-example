import { FC } from 'react';
import Background from './Background';
import { ColorsProvider } from './colors/colors';
import { ColorForm } from './Form';
import { WebsocketContextProvider } from './websocket';

const App: FC = () => {
  return (
    <ColorsProvider>
    <WebsocketContextProvider>
      <ColorForm />
      <Background />
    </WebsocketContextProvider>
    </ColorsProvider>
  );
};

export default App;
