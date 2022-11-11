class Profil extends React.Component {
    constructor(props) {
        super(props);
    }
    async fetchOwnData(mail) {
        return fetch('/api/account/' + mail, {
            method: "GET",
            headers: {
                'Content-Type': 'application/json',
                'token': this.props.token
            }
        }).then(data => data.json())
    }
    async componentDidMount() {
        if (!this.props.role) { //Wenn keine Rolle vorgegeben ist, ist eine Änderung, sonst anlage
            const data = await this.fetchOwnData(this.props.email);
            if (!data) {
                this.props.setMessage("Fehler beim laden der Daten");
            } else {
                this.setState({
                    fetchedData: data,
                    mode: "change"
                });
            }
        } else {
            this.setState({
                fetchedData: {
                    email: "",
                    name: "Neuer Benutzer",
                    role: this.props.role,
                },
                mode: "create"
            })
        }
    }
    async onClick(e) {
        e.preventDefault();
        var update = undefined
        if (this.state.newPw !== "" || this.state.newName !== "") {
            //this.setState({ state: "Daten ans Backend senden ..." })
            this.props.setMessage("Daten ans Backend senden ...");
            update = await this.updatePassword();
        } else {
            this.props.setMessage("Keine Änderung der Daten");
        }

        if (update && update.status < 300) {
            this.props.setMessage("Daten erfolgreich geändert");
        } else {
            this.props.setMessage("Daten konnten nicht geändert werden");
        }
        this.setState({ buttonActive: true })
    }
    async updatePassword() {
        return fetch('/api/account', {
            method: this.state.mode === "create"? "POST": "PUT",
            headers: {
                'Content-Type': 'application/json',
                'token': this.props.token
            },

            body: JSON.stringify({
                "email": this.state.mode === "change" ? this.state.fetchedData.email : this.state.newMail,
                "password": this.state.newPw,
                "role": this.state.fetchedData.role,
                "name": this.state.newName
            })
        })
    }
    render() {
        if (!this.state || !this.state.fetchedData) {
            return (<div />)
        }

        return (
            <div>
                <div className="row">
                    <div className="col-2" />
                    <div className="col-8">
                        <h1>Herzlich Willkommen {this.state.fetchedData.name}</h1>
                    </div>
                </div>

                <div className="row">
                    <div className="col-2" />
                    <div className="col-4">Benachrichtungsadresse:</div>
                    <div className="col-4">

                        {this.state.mode === "change" &&
                            this.state.fetchedData.email
                        }
                        {this.state.mode === "create" &&
                            <input type="text" placeholder={this.state.fetchedData.email} onChange={e => this.setState({ "newMail": e.target.value })} />
                        }
                    </div>
                </div>
                <div className="row">
                    <div className="col-2" />
                    <div className="col-4">Anzeigename</div>
                    <div className="col-4">
                        <input type="text" placeholder={this.state.fetchedData.name} onChange={e => this.setState({ "newName": e.target.value })} />
                    </div>

                </div>

                {this.state.fetchedData.role === "Admin" &&
                    <div className="row">
                        <div className="col-2" />
                        <div className="col-4">Rolle</div>
                        <div className="col-4">{this.state.fetchedData.role}</div>
                    </div>
                }
                <div className="row">
                    <div className="col-2" />
                    <div className="col-4">neues Passwort setzen</div>
                    <div className="col-4"><input type="text" onChange={e => this.setState({ "newPw": e.target.value })} /></div>
                </div>
                <div className="row">
                    <div className="col-2" />
                    <div className="col-8"><button onClick={e => this.onClick(e)}>Sichern</button></div>
                </div>
            </div>
        )

    }
}