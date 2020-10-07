import React, {Component} from 'react';
import {Button,
  Badge,
  ListGroup,
  Container, 
  Row,
  Col, 
  InputGroup, 
  FormControl} from 'react-bootstrap';

import {Redirect} from 'react-router-dom';
import Header from './Header';

class Chat extends Component {

  constructor(props) {

    super(props);

    this.state = {
      ws: null,
      redirect: false,
      message: '',
      messages: [],
      to: '',
      onlineUsers: [],
      disabledLogout: false
    };

    this._handleChange = this._handleChange.bind(this);
    this._handleSendMessage = this._handleSendMessage.bind(this);
    this._handleKeyEnter = this._handleKeyEnter.bind(this);
  }

  componentDidMount() {
    const username = localStorage.getItem('username');
    if (username == null) {
      this.setState({redirect: true});
    } else {
      this.setState({username: username});
      this._connect();
    }
    
  }

  componentWillUnmount() {
    console.log('componen will unmount');
    const {ws} = this.state;
    if (ws != null) {
      ws.close();
    }
    this.setState({disabledLogout: true});
  }

  _connect() {
    const username = localStorage.getItem('username');
    const token = localStorage.getItem('token');

    let ws = new WebSocket(`${this.props.socketBaseUrl}/users/chat`);

    ws.onopen = () => {
        console.log("socket opened..");

        // when success upgrade connection
        // then send auth and user information
        let msg = {
            messageType: 'authMessage',
            authKey: this.props.authKey,
            username: username,
            token: token
        }

        ws.send(JSON.stringify(msg));

        this.setState({ws: ws});
    }

    ws.onclose = (e) => {
        console.log('connection closed.. ', e);
    }

    ws.onerror = (e) => {
        console.log('connection error.. ', e);
        ws.close();
    }

    ws.onmessage = (e) => {
        let messageData = e.data;
        let message = JSON.parse(messageData);

        switch (message.messageType) {
          case 'usersStatus':
            let {onlineUsers} = this.state;
            onlineUsers = message.onlineUsers;

            this.setState({onlineUsers: onlineUsers});
            break;
          case 'authFail':
            // remove localstorge
            localStorage.removeItem('username');
            localStorage.removeItem('token');
            localStorage.removeItem('expired');
            console.log('-----------');
            console.log(message);
            // set redirect to true
            this.setState({redirect: true});
            break;
          default:
            let {messages} = this.state;
            messages.push(message);

            this.setState({messages: messages});
        }
    }
  }

  _handleChange(e) {
    const target = e.target;
    const value = target.value;
    const name = target.name;

    this.setState({[name]: value});
  }

  _handleSendMessage() {
    this._sendMessage();
  }

  _handleKeyEnter(e) {
    const key = e.key;
    if (key === 'Enter') {
      this._sendMessage();
    }
  }

  _sendMessage() {
    const {ws, message, to} = this.state;
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
      content: message
    };

    if (message !== '') {
      
      ws.send(JSON.stringify(msg));
      this.setState({message: ''});
    }
  }

  _renderColumnPeople(props) {
    return (
      <Col sm={4}>
        <h3>People</h3>
        <ListGroup variant="flush">
          {
            props.onlineUsers.map((user, index) => {
              return <ListGroup.Item key={index}>{(user.status ? <Badge pill variant="success">.</Badge> : <Badge pill variant="dark">.</Badge>)} {user.username}</ListGroup.Item>
            })
          }
        </ListGroup>
      </Col>
    );
  }

  _renderColumnMessage(props) {
    return (
      <Col sm={8}>
        <h3>Messages</h3>
        <InputGroup className="mb-3">
          <FormControl name="to" placeholder="to" aria-label="to" value={props.to} onChange={props._handleChange}/>
          <FormControl name="message" placeholder="message" aria-label="message" 
            value={props.message} onChange={props._handleChange} onKeyPress={props._handleKeyEnter}/>
          <InputGroup.Append>
            <Button variant="outline-secondary" onClick={props._handleSendMessage} disabled={!props.message}>Send</Button>
          </InputGroup.Append>
        </InputGroup>
        <ListGroup variant="flush">
          {
            props.messages.map((message, index) => {
              return <ListGroup.Item key={index}>{message.from}: {message.content}</ListGroup.Item>
            })
          }
        </ListGroup>
      </Col>
    );
  }

  render() {
    const {redirect} = this.state;

    if (redirect) {
      return (
        <Redirect to='/'/>
      );
    }

    return (
      <div>
        <Header disabledLogout={this.state.disabledLogout}/>
        <Container>
          <Row>
            {this._renderColumnPeople({onlineUsers: this.state.onlineUsers})}
            {this._renderColumnMessage({
              to: this.state.to,
              _handleChange: this._handleChange,
              _handleKeyEnter: this._handleKeyEnter,
              _handleSendMessage: this._handleSendMessage,
              message: this.state.message,
              messages: this.state.messages,

            })}
            
          </Row>
        </Container>
      </div>
    );
  }
}

export default Chat;
