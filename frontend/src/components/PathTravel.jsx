import path from "path-browserify";

function defaultClick({ absolutePath }) {
  console.log(absolutePath);
}

const itemStyle = {
  display: "inline-block",
  margin: "5px",
  fontSize: "1.5em",
  cursor: "pointer",
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
            onClick={() => {
              onClick({ absolutePath: tmp, path: item.path });
            }}
          >
            {item?.icon}
            {item?.title || item.path}
          </span>
          <span>{splitor}</span>
        </span>
      );
    });
    return result;
  };

  return <div>{genItems(items)}</div>;
}

export default PathTravel;
