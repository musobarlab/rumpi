import React, {Component} from 'react';
import {Chat, Login} from './component';
import {Switch, Route, BrowserRouter} from 'react-router-dom';

class App extends Component {
  render() {
    return (
      <BrowserRouter>
        <div>
          <Switch>
            <Route exact path='/'>
              <Login/>
            </Route>
            <Route path='/chat'>
              <Chat socketUrl={this.props.socketUrl} authKey={this.props.authKey}/>
            </Route>
          </Switch>
        </div>

      </BrowserRouter>
    );
  }
}

App.defaultProps = {
  socketUrl: 'ws://localhost:9000/ws'
};

export default App;
