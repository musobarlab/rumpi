import React, {Component} from 'react';
import {Button,
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
    let ws = new WebSocket(this.props.socketUrl);

    ws.onopen = () => {
        console.log("socket opened..");

        // when success upgrade connection
        // then send auth and user information
        let msg = {
            username: username,
            messageType: 'authMessage',
            authKey: this.props.authKey
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
        if (message.messageType === 'usersStatus') {
          let {onlineUsers} = this.state;
          console.log(message.onlineUsers);
          onlineUsers = message.onlineUsers;
          this.setState({onlineUsers: onlineUsers});
        } else {
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
            <Col sm={4}>
              <h3>People</h3>
              <ListGroup variant="flush">
                {
                  this.state.onlineUsers.map((user, index) => {
                  return <ListGroup.Item key={index}>{user.username}: {(user.status ? 'online' : 'offline')}</ListGroup.Item>
                  })
                }
              </ListGroup>
            </Col>
            <Col sm={8}>
              <h3>Messages</h3>
              <InputGroup className="mb-3">
                <FormControl name="to" placeholder="to" aria-label="to" value={this.state.to} onChange={this._handleChange}/>
                <FormControl name="message" placeholder="message" aria-label="message" 
                  value={this.state.message} onChange={this._handleChange} onKeyPress={this._handleKeyEnter}/>
                <InputGroup.Append>
                  <Button variant="outline-secondary" onClick={this._handleSendMessage} disabled={!this.state.message}>Send</Button>
                </InputGroup.Append>
              </InputGroup>
              <ListGroup variant="flush">
                {
                  this.state.messages.map((message, index) => {
                  return <ListGroup.Item key={index}>{message.from}: {message.content}</ListGroup.Item>
                  })
                }
              </ListGroup>
            </Col>
          </Row>
        </Container>
      </div>
    );
  }
}

export default Chat;
