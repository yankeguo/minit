pub fn str_to_bool(str: &String) -> bool {
    return str.starts_with("t")
        || str.starts_with("T")
        || str.starts_with("y")
        || str.starts_with("Y")
        || str.starts_with("1");
}
