"use client";

import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
  CardContent,
  CardFooter,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";

import api from '@/lib/axios'
import Loader from "@/components/loader";
import { useRouter } from "next/navigation";
import { APIError } from "@/types/APIError";

const LoginPage = () => {
  const { register, handleSubmit } = useForm();

  const router = useRouter();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  // Sign In Handler
  const onLogIn = async () => {
    setIsLoading(true);
    if (email == "") {
      toast.error("Please fill the Email ID Field");
      setIsLoading(false);
      return;
    } else if (password == "") {
      toast.error("Please fill the Password Field");
      setIsLoading(false);
      return;
    }
    try {
      const res = await api.post(
        "/auth/login",
        { email, password },
        { headers: { "Content-Type": "application/json" }, withCredentials: true }
      );
      console.log(res);
      toast.success(res?.data?.message || "Logged In")
      router.push("/dashboard");
    } catch (error) {
      toast.error("Login failed");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex flex-col items-center justify-center h-screen w-screen">
      <p className="text-6xl mb-6">FileVault</p>
      <Card className="w-96">
        <CardHeader className="flex flex-col items-center gap-y-2">
          <CardTitle className="text-2xl">Welcome to FileVault</CardTitle>
          <CardDescription>Sign In with your credentials</CardDescription>
        </CardHeader>
        <CardContent>
          <form className="space-y-4" onSubmit={handleSubmit(onLogIn)}>
            <div className="space-y-2">
              <Label htmlFor="email">Email Id</Label>
              <Input
                type="email"
                id="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">Password</Label>
              <Input
                type="password"
                id="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>
            <CardFooter className="flex-col justify-between p-0 pt-4">
              <Button type="submit">
                {isLoading ? <Loader /> : "Sign In"}
              </Button>
              <div className="mt-6 text-center text-sm">
                Don&apos;t have an account?{" "}
                <a href="/signup" className="underline underline-offset-4 ">
                  Sign up
                </a>
              </div>
            </CardFooter>
          </form>
        </CardContent>
      </Card>
    </div>
  );
};

export default LoginPage;

