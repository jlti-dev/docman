class Benutzer extends React.Component {
    constructor(props) {
        super(props);
    }
    async fetchOwnData() {
        return fetch('/api/accounts', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'token': this.props.token
            }
        }).then(data => data.json())
    }
    async componentDidMount() {
        const data = await this.fetchOwnData();
        if (!data  || data.status > 300) {
            //this.setState({ error: "Fehler beim laden der Daten" });
            this.props.setMessage("Fehler beim laden der Daten")
        } else {
            this.setState({ fetchedData: data});
        }
    }
    onChange(e, data) {
        e.preventDefault();
        this.props.setMessage("Ändern eines fremden Benutzers!")
        this.props.navTo("profile", {email: data.email})
    }
    onCreate(e, data) {
        e.preventDefault();
        this.props.setMessage("Neu Anlage eines Benutzers mit Rolle" + data.role);
        this.props.navTo("profile", data);
    }
    onLink(e, data) {
        e.preventDefault();
        
        this.props.setMessage("Verlinkungen zu Benutzer " + data.name);
        this.props.navTo("link", data);
    }
    render() {
        if (!this.state || !this.state.fetchedData) {
            //this.props.setMessage("Noch keine Daten geladen");
            //return (<div>Noch keine Daten vorhanden</div>)
        } else if (this.state.error) {
            //return (<h1> Fehler beim Laden der Daten</h1>)
        }
        return (
            <div>
                <div className="row">
                    <div className="col-2" />
                    <div className="col-2">Mail-Adresse</div>
                    <div className="col-2">Name</div>
                    <div className="col-1">Dateien</div>
                    <div className="col-1">Passwort</div>
                    <div className="col-1">Benutzerrolle</div>
                    <div className="col-2">letztes Login</div>
                </div>
                {this.state && this.state.fetchedData && this.state.fetchedData.map((data) => {
                    return (
                        <div key={data.ID} className="row">
                            <div className="col-2" />
                            <div className="col-2">{data.email}</div>
                            <div className="col-2">{data.name}</div>
                            <div className="col-1"><button onClick={e => this.onLink(e, data)}>Link</button></div>
                            <div className="col-1"><button onClick={e => this.onChange(e, data)}>Ändern</button></div>
                            <div className="col-1">{data.role}</div>
                            <div className="col-1">{new Date(data.login).toLocaleString()}</div>
                        </div>
                    )
                })
                }
                <div className="row">
                    <div className="col-2" />
                    <div className="col-4">Neuen Benutzer anlegen</div>
                    <div className="col-2"><button onClick={e => this.onCreate(e, {role: "User"})}>Anwender</button></div>
                    <div className="col-2"><button onClick={e => this.onCreate(e, {role: "Admin"})}>Admin</button></div>
                </div>
            </div>
        )
    }
}