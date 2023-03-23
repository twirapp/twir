function convert_variable_name(varName) {
    let withID = varName.replace('_id', '_ID');
    let parts = withID.split('_');
    let build = '';

    parts.forEach(function (part) {
        build += part[0].toUpperCase();
        build += part.substr(1);
    });

    return build;
}

function convert_type_name(typeName) {
    if (typeName === 'string') {
        return 'string';
    } else if (typeName === 'boolean' || typeName === 'Boolean') {
        return 'bool';
    } else if (typeName === 'integer' || typeName === 'int') {
        return 'int';
    } else {
        return '"' + typeName + '"';
    }
}

function convert_table(t, structName) {
    let text = 'type EventSub' + structName + ' struct {\n';

    for (let i = 1; i < t.rows.length; i++) {
        let currRow = t.rows[i];
        let varName = currRow.cells[0].textContent;
        let varType = currRow.cells[1].textContent;
        let varDescription = currRow.cells[2].textContent;

        text += '\t// ' + varDescription + '\n';
        text += '\t' + convert_variable_name(varName) + ' ' + convert_type_name(varType) + ' `json:"' + varName + '"`\n';
    }

    text += '}\n';

    return text;
}

function convert_and_print(structName) {
    // use $0 from chrome developer tools :)
    console.log(convert_table($0, structName));
}
