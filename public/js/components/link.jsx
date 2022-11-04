class Link extends React.Component {
    constructor(props) {
        super(props);
    }
    async fetchMailFiles(mail) {
        let url = "/api"
        if (mail) {
            url += "/link/" + mail;
        } else {
            url += "/file"
        }
        return fetch(url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'token': this.props.token
            }
        }).then(data => data.json())
    }
    async componentDidMount() {
        const files = await this.fetchMailFiles();
        if (!files || files.status > 300) {
            this.props.setMessage("Fehler beim laden der verfügbaren Dateien");
            return
        }
        const links = await this.fetchMailFiles(this.props.email);
        if (!links || links.status > 300) {
            this.props.setMessage("Fehler beim laden der verlinkten Dateien");
            return
        }

        let data = []
        files.forEach(e => {
            var item = {
                logicalname: "",
                physicalname: e.filename,
                assignedto: this.props.email,
                linked: false,
                index: data.length
            }
            links.forEach(l => {
                if (l.physicalname === item.physicalname) {
                    item.logicalname = l.filename;
                    item.linked = true;
                }
            })
            data.push(item);
        });
        debugger;
        this.setState({
            fetchedData: data,
        })
    }
    async onLink(event, data) {
        event.preventDefault();
        this.props.setMessage("Daten werden ans Backend gesendet");
        
        data.linked = !data.linked;
        
        const retVal = await this.linkBackend({
            logicalname: data.logicalname,
            physicalname: data.physicalname,
            assignedto: data.assignedto
        }, data.linked)
        if (!retVal || retVal.status > 300) {
            this.props.setMessage("Fehler beim bearbeiten der Links");
            return
        }
        this.props.setMessage("Daten wurden vom Backend verarbeitet");

        let arr = this.state.fetchedData;
        arr[data.index] = data;
        
        this.setState({fetchedData: arr});
    }
    async linkBackend(data, link) {
        return fetch('/api/file', {
            //PUT = Link anlegen
            //Delete = Link löschen
            method: link ? "PUT" : "DELETE",
            headers: {
                'Content-Type': 'application/json',
                'token': this.props.token
            },
            body: JSON.stringify(data)
        })
    }
    render() {
        return (
            <div>
                <div className="row">
                    <div className="col-2" />
                    <div className="col-3">Physischer Name (Upload Name)</div>
                    <div className="col-3">Logischer Name (Anzeige Name)</div>
                    <div className="col-1">Verlinken</div>
                    <div className="col-1">Link aufheben</div>

                </div>
                {this.state && this.state.fetchedData && this.state.fetchedData.map((data) => {
                    return (
                        <div key={data.physicalname} className="row">
                            <div className="col-2" />
                            <div className="col-3">{data.physicalname}</div>
                            <div className="col-3">
                                {data.linked === true &&
                                    data.logicalname
                                }
                                {data.linked === false &&
                                    <input type="text" onChange={e => data.logicalname = e.target.value} />
                                }
                            </div>
                            <div className="col-1">
                                {data.linked === false &&
                                    <button onClick={(e) => this.onLink(e, data)}>
                                        Link
                                    </button>
                                }
                            </div>
                            <div className="col-1">
                                {data.linked === true &&
                                    <button onClick={(e) => this.onLink(e, data)}>
                                        Unlink
                                    </button>
                                }
                            </div>
                        </div>
                    )
                })}
            </div>
        )
    }
}