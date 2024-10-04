import { FC, useContext } from "react";
import { ColorsContext } from "./colors/colors";

const Background: FC = () => {
  const {currentColor} = useContext(ColorsContext);
  return (
    <div style={{position: 'fixed', top:0, left:0, right:0, bottom: 0, backgroundColor: currentColor, zIndex :-1, transition: 'all 0.4s'}}></div>
  );
}

export default Background;