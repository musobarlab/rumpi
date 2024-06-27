import React, { useState, useEffect } from 'react';
import { Button, Container, Form } from 'react-bootstrap';
import axios from 'axios';
import { Redirect } from 'react-router-dom';
import Header from '../components/Header';

const Register = (props) => {
  const [fullName, setFullName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [repassword, setRepassword] = useState('');
  const [redirect, setRedirect] = useState(false);
  const [disabledLogout, setDisabledLogout] = useState(false);
  const [disabledSubmit, setDisabledSubmit] = useState(true);
  const [messageRegister, setMessageRegister] = useState('');
  const [messageConfirmPassword, setMessageConfirmPassword] = useState('');

  useEffect(() => {
    const username = localStorage.getItem('username');
    if (username == null) {
      setDisabledLogout(true);
    }
  }, []);

  const _handleSubmit = (e) => {
    _registerService();
    e.preventDefault();
  };

  const _registerService = async () => {
    try {
      let response = await axios({
        url: `${props.apiBaseUrl}/users/register`,
        method: 'POST',
        headers: {
          'content-type': 'application/json',
        },
        data: {
          fullName: fullName,
          email: email,
          password: password,
        },
        auth: {
          username: 'user',
          password: '123456',
        },
      });

      let data = response.data;

      if (response.status !== 201) {
        setMessageRegister(data.message);
      } else {
        setFullName('');
        setEmail('');
        setPassword('');
        setRepassword('');
        setRedirect(true);
      }
    } catch (err) {
      let data = err.response;
      let message = 'please try again later';
      if (data) {
        message = data.data.message;
      }
      setFullName('');
      setEmail('');
      setPassword('');
      setRepassword('');
      setMessageRegister(message);
    }
  };

  const _handleChange = (e) => {
    const { name, value } = e.target;
    if (name === 'fullName') setFullName(value);
    else if (name === 'email') setEmail(value);
    else if (name === 'password') setPassword(value);
    else if (name === 'repassword') setRepassword(value);
  };

  const _handlePasswordConfirm = (e) => {
    const repassword = e.target.value;

    if (password !== repassword) {
      setMessageConfirmPassword('password did not match');
      setDisabledSubmit(true);
    } else {
      setMessageConfirmPassword('');
      setDisabledSubmit(false);
    }
  };

  if (redirect) {
    return <Redirect to='/' />;
  }

  return (
    <div>
      <Header disabledLogout={disabledLogout} />
      <Container>
        <Form onSubmit={_handleSubmit}>
          <Form.Text className="text-muted">{messageRegister}</Form.Text>

          <Form.Group controlId="formFullname">
            <Form.Label>Fullname</Form.Label>
            <Form.Control
              name="fullName"
              type="text"
              placeholder="Full Name"
              value={fullName}
              onChange={_handleChange}
              required
            />
          </Form.Group>

          <Form.Group controlId="formEmail">
            <Form.Label>Email address</Form.Label>
            <Form.Control
              name="email"
              type="email"
              placeholder="Enter email"
              value={email}
              onChange={_handleChange}
              required
            />
          </Form.Group>

          <Form.Group controlId="formPassword">
            <Form.Label>Password</Form.Label>
            <Form.Control
              name="password"
              type="password"
              placeholder="Password"
              value={password}
              onChange={_handleChange}
              required
            />
          </Form.Group>

          <Form.Text className="text-muted">{messageConfirmPassword}</Form.Text>
          <Form.Group controlId="formConfirmPassword">
            <Form.Label>Confirm Password</Form.Label>
            <Form.Control
              name="repassword"
              type="password"
              onKeyUp={_handlePasswordConfirm}
              placeholder="Password Confirmation"
              required
            />
          </Form.Group>

          <Button variant="outline-secondary" type="submit" disabled={disabledSubmit}>
            Register
          </Button>
        </Form>
      </Container>
    </div>
  );
};

export default Register;