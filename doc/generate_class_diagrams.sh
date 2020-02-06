#!/usr/bin/env bash

modules="storage server logging"

pushd ..
for module in $modules
do
    goplantuml -recursive ${module}/ > doc/class_diagram_${module}.uml
    java -jar ~/tools/plantuml.jar doc/class_diagram_${module}.uml
done

popd
