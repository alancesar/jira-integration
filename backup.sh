destination_path=${1:-.}
datetime=$(date +"%Y%m%d%H%M%S")
filename="${datetime}.zip"
temp="tmp-${datetime}"

mkdir -p "$temp"
cp "sqlite.db" "$temp"
docker cp metabase:/metabase.db/metabase.db.mv.db "$temp"
zip -rj "$temp/$filename" "$temp"
mv "$temp/$filename" "$destination_path"
rm -rf "$temp"
