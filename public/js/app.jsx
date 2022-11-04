"use-strict";

class App extends React.Component {
    constructor(props) {
        super(props);
        this.state = {};
    }
    setToken(token) {
        if (token) {
            this.setState(token);
        } else {
            this.setState({
                token: undefined,
                mail: undefined
            });
        }
    }
    setMessage(message) {
        this.setState({ message: message });
    }
    setNavigation(navTo, obj) {
        this.setState({ navTo: navTo, obj: obj })
    }
    render() {
        if (!this.state || !this.state.token) {
            return <Login setToken={this.setToken.bind(this)} 
                            navTo={this.setNavigation.bind(this)} />
        }

        return (
            <div>
                <Header app={this} />
                <div>

                    {this.state.navTo === "profile" &&
                        <Profil setMessage={this.setMessage.bind(this)}
                            email={this.state.obj ? this.state.obj.email : this.state.email}
                            role={this.state.obj ? this.state.obj.role : undefined}
                            token={this.state.token} />}
                    {this.state.navTo === "user" &&
                        <Benutzer setMessage={this.setMessage.bind(this)}
                            navTo={this.setNavigation.bind(this)}
                            email={this.state.email}
                            token={this.state.token} />}
                    {this.state.navTo === "download" &&
                        <Download setMessage={this.setMessage.bind(this)}
                            token={this.state.token} />}
                    {this.state.navTo === "upload" &&
                        <Upload setMessage={this.setMessage.bind(this)}
                            token={this.state.token} />}
                    {this.state.navTo === "link" &&
                        <Link setMessage={this.setMessage.bind(this)}
                            email={this.state.obj.email} 
                            token={this.state.token}/>}

                </div>
                <Footer message={this.state.message} />
            </div>
        )

    }

}

ReactDOM.render(React.createElement(App, {}), document.getElementById("App"));
