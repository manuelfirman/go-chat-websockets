import { HStack, VStack, Box, Text, Spacer, Divider } from "@chakra-ui/react";
import React from "react";

interface IFriend {
  name: string;
  email: string;
}
export const Friend = ({ name, email }: IFriend) => {
  // Attributes
  // Context
  // Methods
  // Component
  return (
    <VStack
      w="200px"
      bg="primary"
      borderRadius={10}
      p={4}
      color="buttonText"
    >
      <HStack w="full">
        <Text fontWeight="bold" color="buttonText">{name}</Text>
        <Spacer />
      </HStack>
      <Divider />
      <HStack w="full">
        <Text fontWeight="bold" color="buttonText">{email}</Text>
        <Spacer />
      </HStack>
    </VStack>
  );
};