import React, {Component} from 'react';
import {Button,
  ListGroup,
  ListGroupItem,
  Jumbotron, 
  Container, 
  Row, 
  Col, 
  InputGroup, 
  FormControl} from 'react-bootstrap';

import {Redirect} from 'react-router-dom';

class Chat extends Component {

  constructor(props) {

    super(props);

    this.state = {
      ws: null,
      redirect: false,
      message: '',
      to: '',
      messages: []
    };

    this._handleChange = this._handleChange.bind(this);
    this._handleSendMessage = this._handleSendMessage.bind(this);
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

  _connect() {
    const username = localStorage.getItem('username');
    let ws = new WebSocket('ws://192.168.100.15:9000/ws');

    ws.onopen = () => {
        console.log("socket opened..");

        // when success upgrade connection
        // then send auth and user information
        let msg = {
            username: username,
            messageType: "authMessage",
            authKey: "555abcd"
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
        let {messages} = this.state;
        messages.push(message);
        this.setState({messages: messages});
    }
  }

  _handleChange(e) {
    const target = e.target;
    const value = target.value;
    const name = target.name;

    this.setState({[name]: value});
  }

  _handleSendMessage() {
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

    ws.send(JSON.stringify(msg));

    this.setState({message: ''});
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
        <Jumbotron fluid>
          <Container>
            <h1>Random Chat Demo</h1>
            <p>
              Random Chat Demo Using Golang and React.
            </p>
          </Container>
        </Jumbotron>
        <Container>
          <Row>
            <Col sm={4}>
            <h3>People</h3>
            </Col>
            <Col sm={8}>
              <h3>Messages</h3>
              <InputGroup className="mb-3">
                <FormControl name="to" placeholder="to" aria-label="to" value={this.state.to} onChange={this._handleChange}/>
                <FormControl name="message" placeholder="message" aria-label="message" value={this.state.message} onChange={this._handleChange}/>
                <InputGroup.Append>
                  <Button variant="outline-secondary" onClick={this._handleSendMessage}>Send</Button>
                </InputGroup.Append>
              </InputGroup>
              <ListGroup variant="flush">
                {
                  this.state.messages.map((message, index) => {
                    console.log(message.from);
                  return <ListGroup.Item>{message.from}: {message.content}</ListGroup.Item>
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
