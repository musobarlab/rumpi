import React, { useState, useEffect } from 'react';
import { Button, Badge, ListGroup, Container, Row, Col, InputGroup, FormControl } from 'react-bootstrap';
import { Redirect } from 'react-router-dom';
import Header from '../components/Header';

const Chat = (props) => {
  const [ws, setWs] = useState(null);
  const [redirect, setRedirect] = useState(false);
  const [message, setMessage] = useState('');
  const [messages, setMessages] = useState([]);
  const [to, setTo] = useState('');
  const [onlineUsers, setOnlineUsers] = useState([]);
  const [disabledLogout, setDisabledLogout] = useState(false);
  const [username, setUsername] = useState(localStorage.getItem('username'));

  useEffect(() => {
    if (username == null) {
      setRedirect(true);
    } else {
      _connect();
    }

    return () => {
      if (ws != null) {
        ws.close();
      }
      setDisabledLogout(true);
    };
  }, [username]);

  const _connect = () => {
    const token = localStorage.getItem('token');
    let ws = new WebSocket(`${props.socketBaseUrl}/users/chat`);

    ws.onopen = () => {
      console.log("socket opened..");
      let msg = {
        messageType: 'authMessage',
        authKey: props.authKey,
        username: username,
        token: token,
      };
      ws.send(JSON.stringify(msg));
      setWs(ws);
    };

    ws.onclose = (e) => {
      console.log('connection closed.. ', e);
    };

    ws.onerror = (e) => {
      console.log('connection error.. ', e);
      ws.close();
    };

    ws.onmessage = (e) => {
      let messageData = e.data;
      let message = JSON.parse(messageData);

      switch (message.messageType) {
        case 'usersStatus':
          setOnlineUsers(message.onlineUsers);
          break;
        case 'authFail':
          localStorage.removeItem('username');
          localStorage.removeItem('token');
          localStorage.removeItem('expired');
          console.log('-----------');
          console.log(message);
          setRedirect(true);
          break;
        default:
          setMessages((prevMessages) => [...prevMessages, message]);
      }
    };
  };

  const _handleChange = (e) => {
    const target = e.target;
    const value = target.value;
    const name = target.name;

    if (name === 'message') {
      setMessage(value);
    } else if (name === 'to') {
      setTo(value);
    }
  };

  const _handleSendMessage = () => {
    _sendMessage();
  };

  const _handleKeyEnter = (e) => {
    const key = e.key;
    if (key === 'Enter') {
      _sendMessage();
    }
  };

  const _sendMessage = () => {
    let toUser = '';
    let messageType = 'broadcast';
    if (to !== '') {
      toUser = to;
      messageType = 'privateMessage';
    }

    const msg = {
      from: "",
      to: toUser,
      messageType: messageType,
      content: message,
    };

    if (message !== '') {
      ws.send(JSON.stringify(msg));
      setMessage('');
    }
  };

  const _renderColumnPeople = (props) => {
    return (
      <Col sm={4}>
        <h3>People</h3>
        <ListGroup variant="flush">
          {props.onlineUsers.map((user, index) => {
            return (
              <ListGroup.Item key={index}>
                {user.status ? <Badge pill variant="success">.</Badge> : <Badge pill variant="dark">.</Badge>}
                {user.username}
              </ListGroup.Item>
            );
          })}
        </ListGroup>
      </Col>
    );
  };

  const _renderColumnMessage = (props) => {
    return (
      <Col sm={8}>
        <h3>Messages</h3>
        <InputGroup className="mb-3">
          <FormControl name="to" placeholder="to" aria-label="to" value={props.to} onChange={props._handleChange} />
          <FormControl
            name="message"
            placeholder="message"
            aria-label="message"
            value={props.message}
            onChange={props._handleChange}
            onKeyPress={props._handleKeyEnter}
          />
          <InputGroup.Append>
            <Button variant="outline-secondary" onClick={props._handleSendMessage} disabled={!props.message}>
              Send
            </Button>
          </InputGroup.Append>
        </InputGroup>
        <ListGroup variant="flush">
          {props.messages.map((message, index) => {
            return (
              <ListGroup.Item key={index}>
                {message.from}: {message.content}
              </ListGroup.Item>
            );
          })}
        </ListGroup>
      </Col>
    );
  };

  if (redirect) {
    return <Redirect to='/' />;
  }

  return (
    <div>
      <Header disabledLogout={disabledLogout} />
      <Container>
        <Row>
          {_renderColumnPeople({ onlineUsers: onlineUsers })}
          {_renderColumnMessage({
            to: to,
            _handleChange: _handleChange,
            _handleKeyEnter: _handleKeyEnter,
            _handleSendMessage: _handleSendMessage,
            message: message,
            messages: messages,
          })}
        </Row>
      </Container>
    </div>
  );
};

export default Chat;
