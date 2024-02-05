import React from 'react'
import { useState } from 'react'

import Title from './Title'
import MessageList from './MessageList'
import TypingIndicator from './TypingIndicator'
import InputMessage from './InputMessage'

import './Style.css';
import RoomInfor from './RoomInfor'

/* ChatBox component - composed of Title, MessageList, TypingIndicator, InputMessage */
const ChatBox = ({ owner, ownerAvatar, room, sendMessage, typing, resetTyping, messages, isTyping, isSending }) => {
    return (
      <div className={"chatApp__conv"}>
        <Title room={room} />
       <MessageList owner={owner} messages={messages} room={room}/>
        <div className={"chatApp__convSendMessage clearfix"}>
          <TypingIndicator owner={owner} isTyping={isTyping} />
          <InputMessage
            isSending={isSending}
            owner={owner}
            ownerAvatar={ownerAvatar}
            sendMessage={sendMessage}
            typing={typing}
            resetTyping={resetTyping}
          />
        </div>
      </div>
    );
  };
  

export default ChatBox