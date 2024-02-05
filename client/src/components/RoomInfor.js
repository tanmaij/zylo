import React from 'react'

import './Style.css'

const RoomInfor = ({name, avatar, desc, address}) => {
  return (
    <div className='chatApp__noMesRoomInfor'>
        <img src={avatar} alt={name} className="chatApp__noMesAvatar" />
        <b>{desc}</b>
        <span>Sống tại <b>{address}</b></span>
    </div>
  )
}

export default RoomInfor