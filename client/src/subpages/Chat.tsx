import React, { useRef, useEffect } from "react";
import { useProvider } from "../context";
import { Message } from "../models/message";
import { VStack } from "@chakra-ui/react";
import { ShowMessage } from "../components/ShowMessage";

interface IChat {
  socket: WebSocket | undefined;
  isGroup: boolean;
  toId: string;
  messages: Message[];
  setMessages: React.Dispatch<React.SetStateAction<Message[]>>;
}

export const Chat = ({ socket, isGroup, toId, messages, setMessages }: IChat) => {
  const { user } = useProvider();
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleGetMessages = async () => {
    if (!user) return;
    const m = await user.GetMessages(toId, isGroup);
    if (m) setMessages(m);
  };

  const GetNewMessageId = () => {
    if (!messages || messages.length === 0) return 1;
    return messages[messages.length - 1].id + 1;
  };

  const handleMessages = async () => {
    if (!socket || !user || !messages) return;
    socket.onmessage = function (e) {
      const json = JSON.parse(e.data);
      setMessages([
        ...messages,
        new Message(GetNewMessageId(), json.user_id, Number(toId), json.message, new Date().toISOString()),
      ]);
    };
  };

  useEffect(() => {
    handleMessages();
  }, []);

  useEffect(() => {
    handleGetMessages();
  }, []);

  if (!user || !messages) return null;

  return (
    <VStack w="full" h="690px" overflowY="scroll">
      {messages.map((m, idx) => (
        <ShowMessage key={idx} sender={m.sender_id} message={m.message} owner={user.id === m.sender_id} />
      ))}
      <div ref={messagesEndRef} />
    </VStack>
  );
};
