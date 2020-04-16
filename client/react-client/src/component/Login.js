import React, {Component} from 'react';
import {Button, 
  Jumbotron, 
  Container, 
  InputGroup, 
  FormControl} from 'react-bootstrap';

import {Redirect} from 'react-router-dom';

class Login extends Component {

  constructor(props) {

    super(props);

    this.state = {
      username: '',
      redirect: false
    };

    this._handleChange = this._handleChange.bind(this);
    this._handleLogin = this._handleLogin.bind(this);
  }

  _handleLogin() {
    const {username} = this.state;
    localStorage.setItem('username', username);
    this.setState({redirect: true});
  }

  _handleChange(e) {
    this.setState({username: e.target.value});
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
        <Jumbotron fluid>
          <Container>
            <h1>Random Chat Demo</h1>
            <p>
              Random Chat Demo Using Golang and React.
            </p>
          </Container>
        </Jumbotron>
        <Container>
            <InputGroup className="mb-3">
                <FormControl id="username" placeholder="username" aria-label="username" value={this.state.username} onChange={this._handleChange}/>
                <InputGroup.Append>
                    <Button variant="outline-secondary" onClick={this._handleLogin}>Login</Button>
                </InputGroup.Append>
            </InputGroup>
        </Container>
      </div>
    );
  }
}

export default Login;
