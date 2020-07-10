echo ""
echo "PRE-REQS"
echo "-------"
echo "operator-sdk"

echo "INSTRUCTIONS"
echo "------------"
echo "1) cd to where you want to create the operator project"
echo "2) define the variables when you run the command in the next step (w/ no spaces)"
echo "3) run: OPERATOR_NAME=<operator-project-name> CRD_NAME=<NameOfCRDInUpperCamelCase> source path/to/this/script/create-golang-operator.sh"
echo "4) validate variables from the output below"
echo "5) answer '1' to proceed, '2' to cancel"
echo ""
echo "VARIABLES"
echo "---------"
echo "OPERATOR_NAME: ${OPERATOR_NAME}"
echo "CRD_NAME" : ${CRD_NAME}
echo ""
echo "continue (type '1' or '2')?"

select yn in "Yes" "No"
case $yn in

    Yes ) 
    # create project
    operator-sdk new $OPERATOR_NAME-operator --repo=github.com/keunlee/$OPERATOR_NAME-operator
    cd $OPERATOR_NAME-operator

    # create crd
    operator-sdk add api --api-version=$OPERATOR_NAME.thekeunster.local/v1alpha1 --kind=$CRD_NAME
    operator-sdk generate k8s
    operator-sdk generate crds

    # add controller
    operator-sdk add controller --api-version=$OPERATOR_NAME.thekeunster.local/v1alpha1 --kind=$CRD_NAME
    break;;

    No ) 
    break;;
esac