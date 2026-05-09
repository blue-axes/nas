import path from "path-browserify";

function defaultClick({ absolutePath }) {
  console.log(absolutePath);
}

const itemStyle = {
  display: "inline-block",
  margin: "0 2px",
  fontSize: "1.15em",
  cursor: "pointer",
  color: "#94a3b8",
  padding: "2px 6px",
  borderRadius: "4px",
  transition: "color 0.2s ease, background 0.2s ease",
};

const separatorStyle = {
  color: "#475569",
  margin: "0 2px",
  fontSize: "1em",
};

function PathTravel({ items = [], splitor = "/", onClick = defaultClick }) {
  let genItems = (items) => {
    let result = [];
    let absolutePath = "";
    items?.map((item) => {
      absolutePath = path.join(absolutePath, item.path);
      let tmp = absolutePath;
      result.push(
        <span key={tmp}>
          <span
            style={itemStyle}
            className="tech-breadcrumb-item"
            onClick={() => {
              onClick({ absolutePath: tmp, path: item.path });
            }}
          >
            {item?.icon}
            {item?.title || item.path}
          </span>
          <span style={separatorStyle}>{splitor}</span>
        </span>,
      );
    });
    if (result.length > 0) {
      result.pop();
    }
    return result;
  };

  return <div>{genItems(items)}</div>;
}

export default PathTravel;