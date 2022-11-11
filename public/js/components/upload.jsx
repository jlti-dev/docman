class Upload extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            selectedFile: {
                name: "Keine Datei ausgewählt!"
            }
        }
    }
    onFileChange(event) {

        // Update the state
        this.setState({ selectedFile: event.target.files[0] });

    }
    async sendFile(formData) {
        
        return fetch('/api/file', {
            method: 'POST',
            headers: {
                //'Content-Type': 'multipart/form-data',
                'token': this.props.token
            },
            body:formData
        }).then(data => data.json())
    }
    async onFileUpload() {
        
        this.props.setMessage("Upload started")
        // Create an object of formData
        const formData = new FormData();

        // Update the formData object
        formData.append(
            "uploadFile",
            this.state.selectedFile
        );

        const data = await this.sendFile(formData);
        if (!data || data.status > 300) {
            this.props.setMessage("Fehler beim laden der Daten")
        } else {
            this.setState({fetchedData: data});
            this.props.setMessage("Upload erfolreich")
        }
    }
    render() {

        return (
            <div>
                <div className="row">
                    <div className="col-4" />
                    <div className="col-4">
                        <input type="file" onChange={e => this.onFileChange(e)} />
                    </div>
                </div>
                <div className="row">
                    <div className="col-4" />
                    <div className="col-4">
                        <button onClick={e => this.onFileUpload(e)}>
                            Upload {this.state.selectedFile.name}!
                        </button>
                    </div>
                </div>
            </div>
        )
    }
}