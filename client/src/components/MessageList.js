import React from 'react'

import MessageItem from './MessageItem'

import './Style.css';
import RoomInfor from './RoomInfor';

/* MessageList component - contains all messages */
const MessageList = ({ room,owner, messages }) => {
    return (
      <div className={"chatApp__convTimeline"}>
        {
        messages.length == 0 ? 
        <RoomInfor name={room.name} avatar={room.avatar} desc={room.desc} address={room.address} /> :
         messages.slice(0).reverse().map((messageItem) => (
          <MessageItem
            key={messageItem.id}
            owner={owner}
            sender={messageItem.sender}
            senderAvatar={messageItem.senderAvatar}
            message={messageItem.message}
          />
        ))
        }
      </div>
    );
  };

export default  MessageList