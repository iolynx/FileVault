import { ModeToggle } from "@/components/mode-toggle";
import Link from "next/link";
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  navigationMenuTriggerStyle,
} from "@/components/ui/navigation-menu";
import { ReactNode } from "react";
import { UserAccountPopover } from "@/components/UserAccountPopover";
import { StorageQuotaBanner } from "@/components/StorageQuotaBanner";
import { getCurrentUser } from "@/lib/auth";
import AuthStoreInitializer from "@/components/AuthStoreInitializer";

const AuthorizedLayout = async ({ children }: { children: ReactNode }) => {
  const user = await getCurrentUser();
  return (
    <div className="flex flex-col items-center h-screen w-full">
      <AuthStoreInitializer user={user} />
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

              {/*Admin only Page */}
              {user?.role === 'admin' && (
                <NavigationMenuItem>
                  <Link href="/admin" passHref>
                    <NavigationMenuLink className={navigationMenuTriggerStyle()}>
                      Admin Panel
                    </NavigationMenuLink>
                  </Link>
                </NavigationMenuItem>
              )}

            </NavigationMenuList>
          </NavigationMenu>
        </nav>
        <div>
          <UserAccountPopover />

          <StorageQuotaBanner user={user} />
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

