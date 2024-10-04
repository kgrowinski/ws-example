import React from "react";

type ColorsContextType = {
  currentColor: string;
  setCurrentColor: (color: string) => void;
}
export const ColorsContext = React.createContext<ColorsContextType>({
  currentColor: '#fff',
  setCurrentColor: () => {},
});

type ColorsProviderProps = {
  children: React.ReactNode;
}

export const ColorsProvider: React.FC<ColorsProviderProps> = ({ children }) => {
  const [currentColor, setCurrentColor] = React.useState('#fff');
  return (
    <ColorsContext.Provider value={{ currentColor, setCurrentColor }}>
      {children}
    </ColorsContext.Provider>
  );
}