import React, {Component} from 'react';
import {Chat, Login} from './component';
import {Switch, Route, BrowserRouter} from 'react-router-dom';

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      disabledLogout: false
    }
    
  }

  componentDidMount() {
    const username = localStorage.getItem('username');
    if (username == null) {
      this.setState({disabledLogout: true});
    }
  }

  render() {
    return (
      <BrowserRouter>
        <div>
          <Switch>
            <Route exact path='/'>
              <Login disabledLogout={this.state.disabledLogout}/>
            </Route>
            <Route path='/chat'>
              <Chat socketUrl={this.props.socketUrl} authKey={this.props.authKey} disabledLogout={this.state.disabledLogout}/>
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
