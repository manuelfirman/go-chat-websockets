import { ArrowForwardIcon } from "@chakra-ui/icons";
import { Button, HStack, Input, Box, VStack } from "@chakra-ui/react";
import React, { useRef, useEffect } from "react";
import { useProvider } from "../context";
import { Message } from "../models/message";

interface IInputMessage {
  socket: WebSocket | undefined;
  isGroup: boolean;
  toId: string;
  messages: Message[];
  setMessages: React.Dispatch<React.SetStateAction<Message[]>>;
}

export const InputMessage = ({ socket, isGroup, toId, messages, setMessages }: IInputMessage) => {
  const [message, setMessage] = React.useState<string>("");
  const { user } = useProvider();
  const inputRef = useRef<HTMLInputElement>(null);

  const handleSendMessage = async () => {
    if (!socket || !user) return;
    const obj = {
      user_id: user.id,
      is_jwt: false,
      is_group: isGroup,
      to_id: Number(toId),
      message: message,
    };
    const jsonString = JSON.stringify(obj);
    socket.send(jsonString);
    setMessages([
      ...messages,
      new Message(
        messages.length + 1,
        user.id,
        Number(toId),
        message,
        new Date().toISOString()
      ),
    ]);
    setMessage("");
  };

  // Auto scroll to bottom when new message is added
  useEffect(() => {
    inputRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  return (
    <VStack w="full" spacing={4} mb={4}>
      <HStack w="full" bg="black" borderRadius={10}>
        <Input
          placeholder="Write your message here..."
          w="100%"
          bg="black"
          value={message}
          onChange={(e) => setMessage(e.currentTarget.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              handleSendMessage();
            }
          }}
        />
        <Button variant="signup" w="50px" onClick={handleSendMessage}>
          <ArrowForwardIcon />
        </Button>
      </HStack>
      <Box ref={inputRef} />
    </VStack>
  );
};
