import React, {Component} from 'react';
import {
  Button, 
  Container,
  Form} from 'react-bootstrap';

import axios from 'axios';
import {Redirect} from 'react-router-dom';
import Header from './Header';

class Login extends Component {

  constructor(props) {

    super(props);

    this.state = {
      username: '',
      password: '',
      redirect: false,
      disabledLogout: false,
      messageLogin: ''
    };

    this._handleChange = this._handleChange.bind(this);
    this._handleSubmit = this._handleSubmit.bind(this);
  }

  componentDidMount() {
    const username = localStorage.getItem('username');
    if (username == null) {
      this.setState({disabledLogout: true});
    }
  }

  _handleSubmit(e) {

    const {username, password} = this.state;

    axios({
      url: 'http://192.168.100.15:9000/users/login',
      method: 'POST',
      headers: { 
        'content-type': 'application/json'
      },
      data: {
        'username': username,
        'password': password
      },
      auth: {
        username: 'user',
        password: '123456'
      }
    }).then((response) => {
      let data = response.data;

      if (response.status !== 200) {
        let {messageLogin} = this.state;
        messageLogin = data.message
        this.setState({messageLogin: messageLogin});
      } else {
        localStorage.setItem('username', data.data.email.split('@')[0]);
        localStorage.setItem('token', data.data.accessToken);
        localStorage.setItem('expired', data.data.accessTokenExpired);
        this.setState({username: '', password: '', redirect: true});
      }
    }).catch((err) => {
      let data = err.response;
      let {messageLogin} = this.state;
      if (data.data) {
        messageLogin = data.data.message;
      } else {
        messageLogin = 'please try again later';
      }

      this.setState({username: '', password: ''});
      this.setState({messageLogin: messageLogin});
    });

    e.preventDefault();
  }

  _handleChange(e) {
    const target = e.target;
    const value = target.value;
    const name = target.name;

    this.setState({[name]: value});
  }

  render() {
    const {redirect} = this.state;

    if (redirect) {
      return (
        <Redirect to='/chat'/>
      );
    }

    return (
      <div>
        <Header disabledLogout={this.state.disabledLogout}/>
        <Container>
          <Form onSubmit={this._handleSubmit}>
            <Form.Text className="text-muted">
            {this.state.messageLogin}
            </Form.Text>
            <Form.Group controlId="formBasicEmail">
              <Form.Label>Email address</Form.Label>
              <Form.Control name="username" type="email" placeholder="Enter email" value={this.state.username} onChange={this._handleChange} required={true}/>
            </Form.Group>

            <Form.Group controlId="formBasicPassword">
              <Form.Label>Password</Form.Label>
              <Form.Control name="password" type="password" placeholder="Password" value={this.state.password} onChange={this._handleChange} required={true}/>
            </Form.Group>
            <Button variant="outline-secondary" type="submit">
              Login
            </Button>
          </Form>
        </Container>
      </div>
    );
  }
}

export default Login;
