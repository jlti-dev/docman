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
    const token = await this.loginUser(
      this.state.user,
      this.state.pw
    );
    
    if (!token) {
      this.setState({
        error: "Login fehlgeschlagen"
      })
    } else {
      this.setState({
        error: ""
      });
      //komplettes token zurückgeben, da information über eigene Mail
      this.props.setToken(token);
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
    //<form>
        
    //    <label>
    //      <p>Username</p>
    //      <input type="text" onChange={e => this.setState({ "user": e.target.value })} />
    //    </label>
    //    <label>
    //      <p>Password</p>
    //      <input type="password" onChange={e => this.setState({ "pw": e.target.value })} />
    //    </label>
    //    {this.state.error && <p>{this.state.error}</p> }
    //    <span />
    //    <div>
    //      <button type="submit" onClick={e => this.onClick(e)}>Submit</button>
    //    </div>
    //  </form>
  }

}