import React, { useState, useEffect } from 'react';
import { Button, 
  Container,
  Form } from 'react-bootstrap';

import axios from 'axios';
import { Redirect } from 'react-router-dom';
import Header from '../components/Header';

export default function Login(props) {

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [redirect, setRedirect] = useState(false);
  const [disabledLogout, setDisabledLogout] = useState(false);
  const [messageLogin, setMessageLogin] = useState('');

  useEffect(() => {
    const username = localStorage.getItem('username');
    if (username == null) {
      setDisabledLogout(true);
    }
  }, [disabledLogout]);

  const handleSubmit = (e) => {
    loginService();
    
    e.preventDefault();
  };

  const loginService = async () => {

    try {
      let response = await axios({
        url: `${props.apiBaseUrl}/users/login`,
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
      });

      let data = response.data;
      if (response.status !== 200) {
        let _messageLogin = data.message
        setMessageLogin(_messageLogin);
      } else {
        localStorage.setItem('username', data.data.email.split('@')[0]);
        localStorage.setItem('token', data.data.accessToken);
        localStorage.setItem('expired', data.data.accessTokenExpired);

        setUsername('');
        setPassword('');
        setRedirect(true);
      }
    } catch(err){
      let data = err.response;
      console.log(err);
      let _messageLogin = 'please try again later';
      if (data) {
        _messageLogin = data.data.message;
      }

      setMessageLogin(_messageLogin);
      setUsername('');
      setPassword('');
    }
  };

  if (redirect) {
    return (
      <Redirect to='/chat'/>
    );
  }

  return (
    <div>
      <Header disabledLogout={disabledLogout}/>
      <Container>
        <Form onSubmit={handleSubmit}>
          <Form.Text className="text-muted">
            {messageLogin}
          </Form.Text>
          <Form.Group controlId="formBasicEmail">
            <Form.Label>Email address</Form.Label>
            <Form.Control name="username" type="email" placeholder="Enter email" value={username} onChange={(e) => setUsername(e.target.value)} required={true}/>
          </Form.Group>

          <Form.Group controlId="formBasicPassword">
            <Form.Label>Password</Form.Label>
            <Form.Control name="password" type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} required={true}/>
          </Form.Group>
          <Button variant="outline-secondary" type="submit">
            Login
          </Button>
        </Form>
      </Container>
    </div>
  );
}
