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

const SignUpPage = () => {
  const { register, handleSubmit } = useForm();

  const router = useRouter();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [name, setName] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  // Sign Up Handler
  const onSignUp = async () => {
    setIsLoading(true);
    if (email == "") {
      toast.error("Please fill the Email ID Field");
      setIsLoading(false);
      return;
    } else if (password == "") {
      toast.error("Please fill the Password Field");
      setIsLoading(false);
      return;
    } else if (name == "") {
      toast.error("Please fill in your Name");
      setIsLoading(false);
      return;
    }
    try {
      const res = await api.post(
        "/auth/signup",
        { email, name, password },
        { headers: { "Content-Type": "application/json" }, withCredentials: true }
      );
      console.log(res);
      toast.success(res?.data?.message || "Account Created")
      router.push("/login");
    } catch (error) {
      toast.error("Sign Up failed");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex flex-col items-center justify-center h-screen w-screen">
      <p className="text-6xl mb-6">FileVault</p>
      <Card className="w-96">
        <CardHeader className="flex flex-col items-center gap-y-2">
          <CardTitle className="text-2xl"> Create An Account </CardTitle>
          <CardDescription>Enter the details below to create your account</CardDescription>
        </CardHeader>
        <CardContent>
          <form className="space-y-3" onSubmit={handleSubmit(onSignUp)}>
            <div className="space-y-2">
              <Label htmlFor="employeeId">Email Id</Label>
              <Input
                type="email"
                id="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="employeeId">Name</Label>
              <Input
                type="text"
                id="name"
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="employeeId">Password</Label>
              <Input
                type="password"
                id="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>
            <CardFooter className="flex-col justify-between p-0 pt-4">
              <Button type="submit">
                {isLoading ? <Loader /> : "Sign Up"}
              </Button>
              <div className="mt-6 text-center text-sm">
                Already have an account?{" "}
                <a href="/login" className="underline underline-offset-4 ">
                  Log In
                </a>
              </div>
            </CardFooter>
          </form>
        </CardContent>
      </Card>
    </div>
  );
};

export default SignUpPage;

