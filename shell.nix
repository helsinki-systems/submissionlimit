with import <nixpkgs> {};

mkShell {
  name = "submissionlimit";
  nativeBuildInputs = [
    go_1_18
    gnumake
  ];

  shellHook = ''
    export GOPATH=$PWD/gopath
  '';
}
