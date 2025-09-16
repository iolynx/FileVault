import { ModeToggle } from "@/components/mode-toggle";
import React, { ReactNode } from "react";

const Layout = ({ children }: { children: ReactNode }) => {
  return (
    <>
      {children}
      <div className="fixed bottom-4 right-4">
        <ModeToggle />
      </div>
    </>
  );
};

export default Layout;
