class Download extends React.Component {
    constructor(props) {
        super(props);
    }
    async fetchOwnData(mail) {
        return fetch('/api/file', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'token': this.props.token
            }
        }).then(data => data.json())
    }
    async componentDidMount() {
        const data = await this.fetchOwnData();
        if (!data) {
            this.setState({ error: "Fehler beim laden der Daten" });
        } else {
            this.setState({ fetchedData: data });
        }
        
    }
    onDownload(filename){
        fetch("/api/file/" + filename, {
            method: "GET",
            headers: {
                "token": this.props.token
            }
        }).then((stream) => {
            return stream.blob();
        })
        .then((data) => {
          var a = document.createElement("a");
          a.href = window.URL.createObjectURL(data);
          a.download = filename;
          a.click();
        }); 
    }
    render(){
        if (!this.state || !this.state.fetchedData) {
            //return (<div className="row">Keine Daten vorhanden</div>)
        } else if (this.state.error) {
            //return (<div className="row"> Fehler beim Laden der Daten</div>)
        }
        
        return(
         <div>
            <div className="row">
                <div className="col-3"/>
                <div className="col-5">
                Hallo, 
                hier folgt eine Erklärung! unpersonalisiert :)
                </div>
            </div>
            {this.state && this.state.fetchedData && this.state.fetchedData.map( (data) => {
                return (
                <div key={data.filename} className="row">
                    <div className="col-3"/>
                    <div className="col-4">{data.filename}</div>
                    <div className="col-1">
                        <button onClick={() => this.onDownload(data.filename)}>
                            Download
                        </button>
                    </div>
                </div>
                )
            })}
         </div>   
        )
    }
}