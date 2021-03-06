import { Router, Route, RouteHandler, Link, browserHistory, IndexRoute } from 'react-router';
import { Nav, Navbar, NavItem, NavDropdown, MenuItem, Button, Input, Glyphicon, Badge } from 'react-bootstrap';
import { Dashboard } from './dashboard.js';
import { Provider } from 'react-redux';
import { createStore, combineReducers, applyMiddleware } from 'redux';
import { syncHistoryWithStore, routerReducer } from 'react-router-redux';
import * as reducers from './reducers';
var React = require('react');
var ReactDOM = require('react-dom');
require("./wave.css");
require('expose?$!expose?jQuery!jquery');
require("bootstrap-webpack");

const store = createStore(
	combineReducers({
		...reducers,
		routing: routerReducer
	})
)
const history = syncHistoryWithStore(browserHistory, store)

var App = (props) => (
	<div>
		<Navbar inverse fixedTop>
			<Navbar.Header>
				<Navbar.Brand>
					<a href="#"><img className="logo" src="wave.svg" /></a>
				</Navbar.Brand>
			</Navbar.Header>
			<Nav>
				<NavItem eventKey={1} href="#">Visualize</NavItem>
				<NavItem eventKey={2} href="#">Report<Badge className="report-count" pullRight>42</Badge></NavItem>
				<NavItem eventKey={3} href="#">Investigate</NavItem>
			</Nav>
			<Nav>
				<Navbar.Form pullLeft>
					<Input type="text" className="form-control navbar-search" placeholder="Search MAC, SSID, Manufacture..." />
				</Navbar.Form>
			</Nav>
			<Nav pullRight>
				<NavDropdown eventKey={4} title="Hayden Parker" id="basic-nav-dropdown" noCaret>
					<MenuItem eventKey={4.1}>Settings<span className="glyphicon glyphicon-cog pull-right settings-icon"></span></MenuItem>
					<MenuItem divider />
					<MenuItem eventKey={4.2}>Logout<span className="glyphicon glyphicon-log-out pull-right logout-icon"></span></MenuItem>
				</NavDropdown>
			</Nav>
		</Navbar>
		{props.children}
	</div>
)

ReactDOM.render(
	<Provider store={store}>
		<Router history={history}>
			<Route path="/" component={App}>
				<IndexRoute component={Dashboard}/>
				<Route path="dashboard" component={Dashboard} />
			</Route>
		</Router>
	</Provider>,
	document.getElementById('content')
)
