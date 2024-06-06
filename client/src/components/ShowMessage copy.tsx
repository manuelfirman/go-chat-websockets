import { HStack, Spacer, Text, Box, VStack, Circle } from "@chakra-ui/react";
import React from "react";

interface IShowMessage {
  owner: boolean;
  message: string;
  sender: number;
}
export const ShowMessage = ({ owner, message, sender }: IShowMessage) => {
  // Attributes
  // Context
  // Methods
  // Component
  if (owner)
    return (
      <HStack w="full">
        <Spacer />
        <HStack bg="outgoingMessageBg" borderRadius={10} p={2}>
          <Text p="5px" color="text">{message}</Text>
        </HStack>
        <Box w="10px" />
      </HStack>
    );
  return (
    <HStack w="full">
      <Box w="10px" />
      <Circle size="40px" bg="secundary">
        <Text color="buttonText">{sender}</Text>
      </Circle>
      <HStack bg="incomingMessageBg" borderRadius={10} p={2}>
        <Text p="5px" color="text">{message}</Text>
      </HStack>
      <Spacer />
    </HStack>
  );
};