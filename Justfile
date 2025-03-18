default:
    just --list

record-demos:
    ls -1 demos/*.tape | xargs -I {} vhs "{}"