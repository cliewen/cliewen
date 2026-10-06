if clue validate > /out/validate.txt 2>&1; then echo "validate=green"; else echo "validate=red"; fi
fm() { [ -f docs/vision.md ] && awk -v k="$1" 'NR==1&&$0!="---"{exit} NR>1&&$0=="---"{exit} index($0,k":")==1{sub(k": *","");print;exit}' docs/vision.md; }
if [ -f docs/vision.md ]; then
  echo "vision-status=$(fm status)"
  echo "vision-provenance=$(fm provenance)"
else echo "vision-status=absent"; echo "vision-provenance=absent"; fi
echo "goals=$(ls docs/goals/G-*.md 2>/dev/null | wc -l)"
git log --stat --format='--- %h %s' trial-base..HEAD > /out/post-gitlog.txt 2>&1
