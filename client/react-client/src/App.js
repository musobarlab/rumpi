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
              <Chat/>
            </Route>
          </Switch>
        </div>

      </BrowserRouter>
    );
  }
}

export default App;
