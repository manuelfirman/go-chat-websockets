import React from "react";
import { HStack, VStack, Box, Text, Spacer, Divider } from "@chakra-ui/react";

interface IGroup {
  name: string;
  description: string;
}
export const Group = ({ name, description }: IGroup) => {
  // Attributes
  // Context
  // Methods
  // Component
  return (
    <VStack
      w="300px"
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
        <Text fontWeight="bold" color="buttonText">{description}</Text>
        <Spacer />
      </HStack>
    </VStack>
  );
};