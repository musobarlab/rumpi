import React, {Component} from 'react';
import {
  Button, 
  Container,
  Form} from 'react-bootstrap';

import axios from 'axios';
import { Redirect } from 'react-router-dom';
import Header from '../components/Header';

class Register extends Component {

  constructor(props) {

    super(props);

    this.state = {
        fullName: '',
        email: '',
        password: '',
        redirect: false,
        disabledLogout: false,
        disabledSubmit: true,
        messageRegister: '',
        messageConfirmPassword: ''
    };

    this._handleChange = this._handleChange.bind(this);
    this._handleSubmit = this._handleSubmit.bind(this);
    this._handlePasswordConfim = this._handlePasswordConfim.bind(this);
  }

  componentDidMount() {
    const username = localStorage.getItem('username');
    if (username == null) {
        this.setState({disabledLogout: true});
    }
  }

  _handleSubmit(e) {
    this._registerService();
    e.preventDefault();
  }

  async _registerService() {
    const {fullName, email, password} = this.state;

    try {
      let response =  await axios({
        url: `${this.props.apiBaseUrl}/users/register`,
        method: 'POST',
        headers: { 
          'content-type': 'application/json'
        },
        data: {
          'fullName': fullName,
          'email': email,
          'password': password
        },
        auth: {
          username: 'user',
          password: '123456'
        }
      });

      let data = response.data;

      if (response.status !== 201) {
        let {messageRegister} = this.state;
        messageRegister = data.message
        this.setState({messageRegister: messageRegister});
      } else {
  
        this.setState({fullName: '', email: '', password: '', repassword: '', redirect: true});
      }

    } catch(err) {
      let data = err.response;
      let {messageRegister} = this.state;
      messageRegister = 'please try again later';
      if (data) {
          messageRegister = data.data.message;
      }

      this.setState({fullName: '', email: '', password: '', repassword: ''});
      this.setState({messageRegister: messageRegister});
    }
  }

  _handleChange(e) {
    const target = e.target;
    const value = target.value;
    const name = target.name;

    this.setState({[name]: value});
  }

  _handlePasswordConfim(e) {
    const target = e.target;
    const repassword = target.value;

    let disabledSubmit = true;
    let {password, messageConfirmPassword} = this.state;
    if (password !== repassword) {
      messageConfirmPassword = 'password did not match';
    } else {
        disabledSubmit = false;
        messageConfirmPassword = '';
    }

    this.setState({messageConfirmPassword: messageConfirmPassword, disabledSubmit: disabledSubmit});
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
          <Form onSubmit={this._handleSubmit}>
            <Form.Text className="text-muted">
                {this.state.messageRegister}
            </Form.Text>

            <Form.Group controlId="formFullname">
              <Form.Label>Fullname</Form.Label>
              <Form.Control name="fullName" type="text" placeholder="Full Name" value={this.state.fullName} onChange={this._handleChange} required={true}/>
            </Form.Group>

            <Form.Group controlId="formEmail">
              <Form.Label>Email address</Form.Label>
              <Form.Control name="email" type="email" placeholder="Enter email" value={this.state.email} onChange={this._handleChange} required={true}/>
            </Form.Group>

            <Form.Group controlId="formPassword">
              <Form.Label>Password</Form.Label>
              <Form.Control name="password" type="password" placeholder="Password" value={this.state.password} onChange={this._handleChange} required={true}/>
            </Form.Group>

            <Form.Text className="text-muted">
                {this.state.messageConfirmPassword}
            </Form.Text>
            <Form.Group controlId="formConfirmPassword">
              <Form.Label>Confirm Password</Form.Label>
              <Form.Control name="repassword" type="password"
                onKeyUp={this._handlePasswordConfim}
                placeholder="Password Confirmation" required={true}/>
            </Form.Group>

            <Button variant="outline-secondary" type="submit" disabled={this.state.disabledSubmit}>
              Register
            </Button>
          </Form>
        </Container>
      </div>
    );
  }
}

export default Register;
