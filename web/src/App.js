import logo from './logo.svg';
import './App.css';
import { useCookies } from 'react-cookie';
import { Component } from 'react';
import axios from 'axios';
axios.defaults.baseURL = process.env.REACT_APP_API_URL;

function App() {
  const [cookies] = useCookies(["is_logged_in"]);

  function isLoggedIn() {
    return cookies.is_logged_in === "true";
  }

  function createNewLink(data) {
    console.log(data);
    return axios.post(`/links/${data.id}`, {"url": data.url}, {withCredentials: true})
  }

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <LinkFrom onCreate={createNewLink} />
        <SignButton isLoggedIn={isLoggedIn()} />
      </header>
    </div>
  );
}

function SignButton(props) {
  const isLoggedIn = props.isLoggedIn;
  if (!isLoggedIn) {
    return <a href={`${axios.defaults.baseURL}/auth/signin`}>로그인</a>
  } else {
    return <a href={`${axios.defaults.baseURL}/auth/signout`}>로그아웃</a>
  }
}

class LinkFrom extends Component {
  state = {
    id: '',
    url: ''
  }

  handleChange = (e) => {
    this.setState({
      [e.target.name]: e.target.value
    })
  }
  handleSubmit = (e) => {
    e.preventDefault();
    this.props.onCreate(this.state)
      .then(response => {
        if (this.response.status === 200) {
          this.setState({
            id: '',
            url: ''
          })
        }
      })
      .catch(error => {
        console.log(this.error);
      });
  }

  render() {
    return (
      <form onSubmit={this.handleSubmit}>
        <input placeholder="ID" value={this.state.id} onChange={this.handleChange} name="id" />
        <input placeholder="URL" value={this.state.url} onChange={this.handleChange} name="url" />
        <input type="submit" value="등록" />
      </form>
    )
  }
}

export default App;
