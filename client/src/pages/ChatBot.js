import React, { useState } from "react";

import ChatRoom from "../components/ChatRoom";
import { detectURL } from "../utils/URL";
import { OWNER_NAME } from "../constants/Key";

const ChatBot = () => {
  const owner = {
    name: OWNER_NAME,
  };

  const [isLoading, setIsLoading] = useState(false);

  const [room, setRoom] = useState({
    name: "Con",
    avatar: "https://i.pravatar.cc/150?img=32",
    desc: "Đang hoạt động",
    address: "Thành phố Hồ Chí Minh",
  });

  const [messages, setMessages] = useState([]);

  const [isTyping, setIsTyping] = useState({
    ricon:true
  });
  const [isSending, setIsSending] = useState(false);

  const sendMessage = (sender, senderAvatar, message) => {
    setIsSending(true)
    setTimeout(() => {
      const messageFormat = detectURL(message);
      const newMessageItem = {
        id: messages.length + 1,
        sender: sender,
        senderAvatar: senderAvatar,
        message: messageFormat,
      };

      setIsSending(false)
      setMessages([...messages, newMessageItem]);
    }, 400);
  };

  const typing = (writer) => {};

  const resetTyping = (writer) => {};

  return (
    <>
      {isLoading ? (
        <section id="chatApp" className="chatApp">
          <div className="chatApp__loaderWrapper">
            <div className="chatApp__loaderText">Loading...</div>
            <div className="chatApp__loader"></div>
          </div>
        </section>
      ) : (
        <ChatRoom
          sendMessage={sendMessage}
          messages={messages}
          typing={typing}
          resetTyping={resetTyping}
          isTyping={isTyping}
          owner={owner}
          room={room}
          isSending={isSending}
        />
      )}
    </>
  );
};

export default ChatBot;
