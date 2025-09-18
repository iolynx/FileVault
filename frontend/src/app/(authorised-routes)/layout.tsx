import { ModeToggle } from "@/components/mode-toggle";
import { Link } from "@radix-ui/react-navigation-menu";
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  navigationMenuTriggerStyle,
} from "@/components/ui/navigation-menu";
import { ReactNode } from "react";

const AuthorizedLayout = ({ children }: { children: ReactNode }) => {
  return (
    <div className="flex flex-col items-center h-screen w-full">
      <div className="flex flex-row w-full place-content-around">
        <nav className="w-full">
          <NavigationMenu className="mx-2">
            <NavigationMenuList className="flex justify-center gap-2 my-2">
              <NavigationMenuItem>
                <NavigationMenuLink
                  asChild
                  className={navigationMenuTriggerStyle()}
                >
                  <Link href="/dashboard">Dashboard</Link>
                </NavigationMenuLink>
              </NavigationMenuItem>
              <NavigationMenuItem>
                <NavigationMenuLink
                  asChild
                  className={navigationMenuTriggerStyle()}
                >
                  <Link href="/storage">Storage</Link>
                </NavigationMenuLink>
              </NavigationMenuItem>
            </NavigationMenuList>
          </NavigationMenu>
        </nav>
        <div>
        </div>
      </div>
      <section className="w-full px-4">{children}</section>
      <div className="fixed bottom-4 right-4">
        <ModeToggle />
      </div>
    </div>
  );
};

export default AuthorizedLayout;

