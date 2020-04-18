import React, {Component} from 'react';
import {
  Button, 
  Container,
  Form} from 'react-bootstrap';

import axios from 'axios';
import {Redirect} from 'react-router-dom';
import Header from './Header';

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
        messageRegister: ''
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

    const {fullName, email, password} = this.state;

    axios({
      url: 'http://192.168.100.15:9000/users/register',
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
    }).then((response) => {
      let data = response.data;

      if (response.status !== 201) {
        let {messageRegister} = this.state;
        messageRegister = data.message
        this.setState({messageRegister: messageRegister});
      } else {
  
        this.setState({fullName: '', email: '', password: '', repassword: '', redirect: true});
      }
    }).catch((err) => {
        let data = err.response;
        let {messageRegister} = this.state;
        if (data.data) {
            messageRegister = data.data.message;
        } else {
            messageRegister = 'please try again later';
        }

        this.setState({fullName: '', email: '', password: '', repassword: ''});
        this.setState({messageRegister: messageRegister});
    });

    e.preventDefault();
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
    let {password, messageRegister} = this.state;
    if (password !== repassword) {
        messageRegister = 'password did not match';
    } else {
        disabledSubmit = false;
        messageRegister = '';
    }

    this.setState({messageRegister: messageRegister, disabledSubmit: disabledSubmit});
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

            <Form.Group controlId="formConfirmPassword">
              <Form.Label>Confirm Password</Form.Label>
              <Form.Control name="repassword" type="password"
                onKeyUp={this._handlePasswordConfim}
                placeholder="Password Confirmation" value={this.state.repassword} required={true}/>
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
