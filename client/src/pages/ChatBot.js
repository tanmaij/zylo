import React, { useState } from "react";

import ChatRoom from "../components/ChatRoom";
import { detectURL } from "../utils/URL";
import { mapCodeBlocksToPre } from "../utils/Content";

import axios from "axios";

const socket = new WebSocket(process.env.REACT_APP_WS_ENDPOINT);
const ChatBot = () => {
  socket.onmessage = (message) => {
    const data = JSON.parse(message.data);
    const msg = {
      event_name: data.event_name,
      data: data.data,
    };

    switch (msg.event_name) {
      case "connection":
        const conData = JSON.parse(msg.data);
        const connectionData = {
          uuid: conData.uuid,
        };

        setOwner({name: connectionData.uuid})

        axios
          .get(
            `${process.env.REACT_APP_API_ENDPOINT}/api/v1/conversations/${connectionData.uuid}`
          )
          .then((res) => {
            setMessages(res.data.messages);
            setRoom(res.data.room);
            setIsLoading(false);
          })
          .catch((err) => console.log(err));
          return;
      case "typing":
        const tdt = JSON.parse(msg.data);
        const typingData = {
          name: tdt.name,
        };

        setIsTyping({
          ...isTyping,
          [typingData.name]: true,
        });

        return;
      case "reset-typing":
        const rtdt = JSON.parse(msg.data);
        const resetTypingData = {
          name: rtdt.name,
        };

        setIsTyping(typing =>{
          delete typing[resetTypingData.name]
          setIsTyping(typing)
        });

        return;

      case "chat":
        const cd = JSON.parse(msg.data);
        const chatData = {
          id:   messages.length+1,
          sender: cd.sender,
          senderAvatar: cd.senderAvatar,
          message: mapCodeBlocksToPre(cd.message),
        };
        setMessages([...messages, chatData]);
      return;

      default:
        return
    }
  };

  const [owner, setOwner] = useState({
    name:""
  })

  const [isLoading, setIsLoading] = useState(true);

  const [room, setRoom] = useState({
    name: "Con",
    avatar: "https://i.pravatar.cc/150?img=32",
    desc: "Đang hoạt động",
    address: "Thành phố Hồ Chí Minh",
  });

  const [messages, setMessages] = useState([]);

  const [isTyping, setIsTyping] = useState({});
  const [isSending, setIsSending] = useState(false);

  const sendMessage = (sender, senderAvatar, message) => {
    const messageFormat = detectURL(message);

    setIsSending(true);
    axios.post(
      `${process.env.REACT_APP_API_ENDPOINT}/api/v1/chat/${sender}`,
      {
        message: messageFormat
      },
      {
        headers: {
          'Content-Type': 'application/json'
        }
      }
    )
    .then((res) => {
      const newMsg = res.data
      newMsg.id=messages.length+1
      newMsg.message=mapCodeBlocksToPre(newMsg.message)

      setMessages([...messages, newMsg]);
      setIsSending(false);
    })
    .catch((err) => {window.location.reload(false);});

    return;
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
