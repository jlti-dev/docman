class Footer extends React.Component {
    render() {
        return (
            <footer>
                {this.props.message &&
                    <div className="row">

                        <div className="col-2" />
                        <div className="col-4">Status</div>
                        <div className="col-4">{this.props.message}</div>

                    </div>
                }
                <div className="row">
                    <div className="col-2" />
                    <div className="col-4">
                        <img src="https://www.systema-projekte.de/wp-content/uploads/2016/09/systema-logo-anthrazit.svg" />
                    </div>
                    <div className="col-2">
                        <a href="https://www.systema-projekte.de">Homepage</a>
                    </div>
                    <div className="col-2">
                        <a href="https://www.systema-projekte.de/datenschutzerklaerung">Datenschutz</a>
                    </div>
                </div>
            </footer>
        );

    }
}