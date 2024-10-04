import { FC, useContext } from "react";
import { SketchPicker } from 'react-color';
import { ColorsContext } from "./colors/colors";
import { WebsocketContext, WS_ACTIONS } from "./websocket";

export const ColorForm: FC =() => {
  const {currentColor} = useContext(ColorsContext);
  const {handleSendWSMessage} = useContext(WebsocketContext)
  const handleChangeComplete = (color) => {
    handleSendWSMessage(WS_ACTIONS.WS_SET_COLOR,color.hex);
  }

  
  return (
    <SketchPicker
        color={ currentColor}
        onChangeComplete={handleChangeComplete }
      />
  );
}
