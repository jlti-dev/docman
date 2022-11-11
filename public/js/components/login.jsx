class Login extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      user: "",
      pw: ""
    }
  }
  async loginUser(email, pw) {
    return fetch('/api/signin', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ "email": email, "password": pw})
    }).then(data => data.json()).catch(() => { })

  }
  async onClick(event) {
    event.preventDefault();
    this.props.setMessage("Login versendet")
    const token = await this.loginUser(
      this.state.user,
      this.state.pw
    );
    
    if (!token) {
      this.props.setMessage("Login fehlgeschlagen")
    } else {
      //komplettes token zurückgeben, da information über eigene Mail
      this.props.setToken(token);
      this.props.setMessage("")
      this.props.navTo("download");
      
    }

  }
  render() {
    return (
      <div>
        <div className="row">
          <div className="col-2"/>
          <div className="col-4">Benutzername</div>
          <div className="col-4">
              <input type="text" onChange={e => this.setState({ "user": e.target.value })} />
          </div>
        </div>
        <div className="row">
        <div className="col-2"/>
          <div className="col-4">Passwort</div>
          <div className="col-4">
            <input type="password" onChange={e => this.setState({ "pw": e.target.value })} />
          </div>
        </div>
        <div className="row">
          <div className="col-2"/>
          <div className="col-8">
          <button type="submit" onClick={e => this.onClick(e)}>Submit</button>
          </div>
        </div>
      </div>
      
    )
  }

}