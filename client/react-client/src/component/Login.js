import React, {Component} from 'react';
import {Button, 
  Container,
  InputGroup, 
  FormControl} from 'react-bootstrap';

import {Redirect} from 'react-router-dom';
import Header from './Header';

class Login extends Component {

  constructor(props) {

    super(props);

    this.state = {
      username: '',
      redirect: false,
      disabledLogout: false
    };

    this._handleChange = this._handleChange.bind(this);
    this._handleLogin = this._handleLogin.bind(this);
  }

  componentDidMount() {
    const username = localStorage.getItem('username');
    if (username == null) {
      this.setState({disabledLogout: true});
    }
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
        <Header disabledLogout={this.state.disabledLogout}/>
        <Container>
            <InputGroup className="mb-3">
                <FormControl id="username" placeholder="username" aria-label="username" value={this.state.username} onChange={this._handleChange}/>
                <InputGroup.Append>
                    <Button variant="outline-secondary" onClick={this._handleLogin} disabled={!this.state.username}>Login</Button>
                </InputGroup.Append>
            </InputGroup>
        </Container>
      </div>
    );
  }
}

export default Login;
