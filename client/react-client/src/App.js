import React, { Component } from 'react';
import { Chat, Login, Register } from './pages';
import { Switch, Route, BrowserRouter } from 'react-router-dom';

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      disabledLogout: false
    }
    
  }

  render() {
    return (
      <BrowserRouter>
        <div>
          <Switch>
            <Route exact path='/'>
              <Login apiBaseUrl={this.props.apiBaseUrl}/>
            </Route>
            <Route exact path='/register'>
              <Register apiBaseUrl={this.props.apiBaseUrl}/>
            </Route>
            <Route path='/chat'>
              <Chat socketBaseUrl={this.props.socketBaseUrl} authKey={this.props.authKey}/>
            </Route>
          </Switch>
        </div>

      </BrowserRouter>
    );
  }
}

App.defaultProps = {
  apiBaseUrl: 'http://localhost:9000',
  socketBaseUrl: 'ws://localhost:9000'
};

export default App;
