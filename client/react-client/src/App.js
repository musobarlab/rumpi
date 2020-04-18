import React, {Component} from 'react';
import {Chat, Login, Register} from './component';
import {Switch, Route, BrowserRouter} from 'react-router-dom';

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
              <Login/>
            </Route>
            <Route exact path='/register'>
              <Register/>
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
